package impl

import (
	"bytes"
	"fmt"
	"sync"
	"text/tabwriter"
	"time"
)

// latest info about puppet node
type Node struct {
	State         bool
	Name          string
	PuppetVersion string
	LastReport    int64
	Errors        bytes.Buffer
}

// describe our collection of puppet environments
type EnvironmentCollection struct {
	sm   sync.Mutex
	Env  map[string]NodeCollection
	Conf *Settings
}

// describe collection of nodes in every environment
type NodeCollection struct {
	sm    sync.Mutex
	Nodes map[string]Node
}

// creates new instance of EnvironmentCollection
func (n EnvironmentCollection) NewEnvironmentCollection() EnvironmentCollection {
	n.sm.Lock()
	defer n.sm.Unlock()
	return EnvironmentCollection{Env: make(map[string]NodeCollection)}
}

// Get nodes of the specified environment
func (n *EnvironmentCollection) GetEnvironmentNodes(name string) *NodeCollection {
	n.sm.Lock()
	nodes, ok := n.Env[name]
	if !ok {
		nodes = NodeCollection{Nodes: make(map[string]Node)}
		n.Env[name] = nodes
	}
	n.sm.Unlock()
	return &nodes
}

// Remove first node occurance
func (n *EnvironmentCollection) RemoveNode(host string) string {
	n.sm.Lock()
	defer n.sm.Unlock()
	for _, nodes := range n.Env {
		_, exists := nodes.Nodes[host]
		if exists {
			delete(nodes.Nodes, host)
			return "OK"
		}
	}
	return "Node: " + host + " not exists"
}

// Process puppet (>=4.0) report
func (n *EnvironmentCollection) ProcessReport(report []byte) bool {
	rep, err := ReportItem{}.FromJson(report)
	if err != nil {
		return false
	}
	// remove node from collection since environment could be changed
	n.RemoveNode(rep.Host)

	nodes := *(n.GetEnvironmentNodes(rep.Environment))
	nodes.sm.Lock()
	defer nodes.sm.Unlock()
	node := Node{Name: rep.Host}
	node.Name = rep.Host
	node.PuppetVersion = rep.PuppetVersion
	// clean up old errors
	node.Errors.Reset()

	// analyzing report status and summ of failures
	node.State = rep.Status != "failed" &&
		rep.Metrics.Events.Failed+rep.Metrics.Resources.Failed+rep.Metrics.Resources.FailedToRestart == 0

	// looking for log message with "err" level
	for _, logI := range rep.Logs {
		if logI.Level == "err" {
			node.State = false
			node.Errors.WriteString(logI.Message + "\n")
		}
	}
	// saving current time
	node.LastReport = time.Now().Unix()
	// saving node
	nodes.Nodes[rep.Host] = node
	return true
}

// Get information about nodes in print form
func (n *EnvironmentCollection) GetInfo() string {
	n.sm.Lock()
	// buffer to store out big string
	var sb bytes.Buffer
	defer n.sm.Unlock()
	max_diff := int64(n.Conf.ControlTime * 60)
	now := time.Now().Unix()

	var state string
	// creating colorized table to print
	w := tabwriter.NewWriter(&sb, 1, 0, 1, ' ', tabwriter.TabIndent|tabwriter.Debug)
	fmt.Fprintf(w, " %v\t %v\t %v\t %v\t\n", "Host", "\033[00mStatus\033[0m", "Last report", "Agent version")
	for _, nodes := range n.Env {
		for host, node := range nodes.Nodes {
			if node.State && now-node.LastReport <= max_diff {
				state = "\033[32mOK\033[0m"
			} else {
				state = "\033[31mError\033[0m"
			}

			fmt.Fprintf(w, " %v\t %v\t %v\t %v\t\n", host, state, time.Unix(node.LastReport, 0).Format("02.01.2006 15:04:05"), node.PuppetVersion)
		}
	}
	fmt.Fprintln(w)
	w.Flush()
	return string(sb.Bytes())
}

// Analize enery node at the moment and return
// OK: if all fine (nagios would be happy)
// List of nodes: if something goes wrong
func (n *EnvironmentCollection) ProcessCollectionState(write_errors bool) string {
	n.sm.Lock()
	defer n.sm.Unlock()
	// buffer to store out big string
	var sb bytes.Buffer
	// max income node report delay
	max_diff := int64(n.Conf.ControlTime * 60)
	// just now, since we using mutex
	now := time.Now().Unix()

	// final puppet state
	overall_state := true

	// analyze and create result
	for env, nodes := range n.Env {
		for host, node := range nodes.Nodes {
			if node.State && now-node.LastReport <= max_diff {
				continue
			}
			overall_state = false
			if write_errors {
				if !node.State {
					sb.WriteString(fmt.Sprintf("\033[1;33m%v\033[0m: failed to implement manifest\n", host))
					sb.WriteString(fmt.Sprintf("Environment: %v\n", env))
					sb.WriteString(fmt.Sprintf("Agent version: %v\n", node.PuppetVersion))
					sb.WriteString(fmt.Sprintf("Last report: %v\n", time.Unix(node.LastReport, 0)))
					sb.Write(node.Errors.Bytes())
					sb.WriteString("\n")
				}
				if now-node.LastReport > max_diff {
					sb.WriteString(fmt.Sprintf("\033[1;33m%v\033[0m: out of sync, no report since \"%v\"\n", host, time.Unix(node.LastReport, 0)))
				}
			} else {
				sb.WriteString(fmt.Sprintf("%v\n", host))
			}
		}
	}

	if overall_state {
		return "OK"
	} else {
		return string(sb.Bytes())
	}
}
