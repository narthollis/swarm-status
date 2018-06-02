"use strict";

async function getCurrent() {
    const response = await fetch('/current');
    return response.json();
}

async function getHistory() {
    const response = await fetch('/history');
    return response.json();
}

class OverallStatus extends React.PureComponent {
    render() {
        return React.createElement(
            'div',
            { className: `alert alert-${this.props.className}` },
            React.createElement("h4", { className: 'alert-heading' }, this.props.text)
        );
    }
}

class StatusDisplay extends React.PureComponent {
    render() {
        const status = this.props.status;

        return React.createElement(
            this.props.element,
            { className: `status bg-${status.ClassName}` },
            status.Value,
            React.createElement('i', { className: `status-icon fas ${status.Icon}`})
        );
    }
}

function *groupHistory(history) {
    let count = 0;
    let className = history[0].ClassName;
    for (const period of history) {
        if (className !== period.ClassName) {
            yield {
                width: count,
                className: className
            };

            count = 1;
            className = period.ClassName;
        } else {
            count = count + 1;
        }
    }

    yield {
        width: count,
        className
    };
}

class History extends React.PureComponent {
    render() {
        const groupings = [...groupHistory(this.props.history)];

        return React.createElement(
            'td',
            { className: 'history' },
            React.createElement(
                'div',
                undefined,
                ...groupings.map(g => React.createElement(
                    'div',
                    { className: `bg-${g.className}`, style: { width: `${g.width}px`} }
                ))
            )
        )
    }
}

class Service extends React.PureComponent {
    render() {
        const service = this.props.service;
        const history = this.props.history;

        return React.createElement(
            'tr',
            undefined,
            React.createElement('td', { className: 'service' }, service.Name),
            React.createElement('td', { className: 'counts' }, service.Running, ' / ', service.Replicas),
            React.createElement(StatusDisplay, { element: 'td', status: service.Status }),
            history && React.createElement(History, { history })
        );
    }
}

class Group extends React.Component {
    constructor(props) {
        super(props);

        this.state = { showBody: false, prevStatus: null };

        this.onClick = this.onClick.bind(this);
    }

    onClick() {
        this.setState(state => ({ showBody: !state.showBody }));
    }

    render() {
        const name = this.props.name;
        const group = this.props.group;
        const history = this.props.history;

        console.log(`${name}-${this.id}`);

        return [
            React.createElement(
                'thead',
                { key: 'thead', className: 'thead-dark group', onClick: this.onClick },
                React.createElement(
                    'tr',
                    undefined,
                    React.createElement('th', { colSpan: 2, className: 'group-name' }, name),
                    React.createElement(StatusDisplay, { element: 'th', status: group.Status }),
                    history && React.createElement(History, { history: history.History })
                )
            ),
            this.state.showBody ?
                React.createElement(
                    'tbody',
                    { key: 'tbody', className: 'table-sm' },
                    group.Services.map(service =>
                        React.createElement(
                            Service,
                            { key: service.ID, service, history: history && history.Services[service.ID] }
                        )
                    )
                )
            : null
        ]
    }
}

Group.getDerivedStateFromProps = function (props, prevState) {
    if (props.group.Status.Value !== prevState.prevStatus) {
        return {
            prevStatus: props.group.Status.Value,
            showBody: props.group.Status.Value !== 'Operational'
        }
    }
};

class Table extends React.PureComponent {
    render() {
        return React.createElement(
            'table',
            { key: 'table', className: 'table table-dark table-striped'},
            this.props.children,
            React.createElement(
                'thead',
                { className: 'thead-light table-sm' },
                React.createElement(
                    'tr',
                    undefined,
                    React.createElement('th', undefined, 'Name'),
                    React.createElement('th', undefined, 'Nodes'),
                    React.createElement('th', undefined, 'Status'),
                    React.createElement('th', undefined, 'History'),
                )
            ),
            React.createElement(
                React.Fragment,
                undefined,
                [...Object.entries(this.props.groups)].map(([groupName, group]) => React.createElement(
                    Group,
                    { key: groupName, name: groupName, group, history: this.props.history[groupName] }
                ))
            )
        );
    }
}

const SECOND = 1000;
const MINUTE = 60 * SECOND;

class Root extends React.PureComponent {
    constructor(props) {
        super(props);

        this.state = {
            current: null,
            history: {}
        };

        this.loadCurrent = this.loadCurrent.bind(this);
        this.loadHistory = this.loadHistory.bind(this);
    }

    componentDidMount() {
        this.loadCurrent();
        setInterval(this.loadCurrent, 10 * SECOND);

        this.loadHistory();
        setInterval(this.loadHistory, 5 * MINUTE);
    }

    loadCurrent() {
        getCurrent().then(current => {
            this.setState(() => ({ current }))
        });
    }

    loadHistory() {
        getHistory().then(history => {
            const historyById = {};

            for (const period of history) {
                for (const [groupName, group] of Object.entries(period.Groups)) {
                    if (historyById[groupName] == null) {
                        historyById[groupName] = {
                            History: [],
                            Services: {},
                        };
                    }

                    historyById[groupName].History.push({
                        timestamp: period.Timestamp,
                        ...group.Status,
                    });

                    for (const service of group.Services) {
                        if (historyById[groupName].Services[service.ID] == null) {
                            historyById[groupName].Services[service.ID] = [];
                        }

                        historyById[groupName].Services[service.ID].push({
                            timestamp: period.Timestamp,
                            ...service.Status,
                        });
                    }
                }
            }

            this.setState(() => ({ history: historyById }));
        })
    }

    render() {
        if (this.state.current == null) {
            return null;
        }

        const current = this.state.current;

        const lastUpdated =luxon.DateTime.fromISO(current.Timestamp).toFormat("ccc d LLL yyyy HH:mm:ss");

        return [
            React.createElement(
                OverallStatus,
                { key: 'overall', className: current.OverallStatus.ClassName, text: current.OverallStatusVerbose }
            ),
            React.createElement(
                Table,
                { key: 'table', groups: current.Groups, history: this.state.history },
                React.createElement(
                    'caption',
                    undefined,
                    `Last Updated ${lastUpdated}`
                )
            )
        ]
    }
}

ReactDOM.render(
    React.createElement(Root),
    document.getElementById('react-root')
);
