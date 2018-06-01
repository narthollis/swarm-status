"use strict";

async function getCurrent() {
    const response = await fetch('/current');
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

class Service extends React.PureComponent {
    render() {
        const service = this.props.service;

        return React.createElement(
            'tr',
            undefined,
            React.createElement('td', { className: 'service' }, service.Name),
            React.createElement('td', { className: 'counts' }, service.Running, ' / ', service.Replicas),
            React.createElement(StatusDisplay, { element: 'td', status: service.Status }),
            React.createElement('td', { className: 'history' })
        );
    }
}

class Group extends React.PureComponent {
    constructor(props) {
        super(props);

        this.onClick = this.onClick.bind(this);
    }

    onClick() {
        this.setState(state => ({ showBody: !state.showBody }));
    }

    render() {
        const name = this.props.name;
        const group = this.props.group;

        return [
            React.createElement(
                'thead',
                { key: 'thead', className: 'thead-dark group', onClick: this.onClick },
                React.createElement(
                    'tr',
                    undefined,
                    React.createElement('th', { colSpan: 2, className: 'group-name' }, name),
                    React.createElement(StatusDisplay, { element: 'th', status: group.Status }),
                    React.createElement('td', { className: 'history' })
                )
            ),
            this.state.showBody ?
                React.createElement(
                    'tbody',
                    { key: 'tbody', className: 'table-sm' },
                    group.Services.map(service => React.createElement(Service, { key: service.ID, service }))
                )
            : null
        ]
    }
}

Group.getDerivedStateFromProps = function (props, prevState) {
    return { showBody: props.group.Status.Value !== 'Operational' };
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
            [...Object.entries(this.props.groups)].map(([groupName, group]) => React.createElement(
                Group,
                { key: groupName, name: groupName, group }
            ))
        );
    }
}

class Root extends React.PureComponent {
    constructor(props) {
        super(props);

        this.state = {
            current: null
        };

        this.loadCurrent = this.loadCurrent.bind(this);
    }

    componentDidMount() {
        this.loadCurrent();
        setInterval(this.loadCurrent, 10000); // 10sec
    }

    loadCurrent() {
        getCurrent().then(current => {
            this.setState(() => ({ current }))
        });
    }

    render() {
        if (this.state.current == null) {
            return null;
        }

        const current = this.state.current;

        const lastUpdated =luxon.DateTime.fromISO(current.Timestamp).toFormat("ccc d LLL yyyy HH:mm:ss")

        return [
            React.createElement(
                OverallStatus,
                { key: 'overall', className: current.OverallStatus.ClassName, text: current.OverallStatusVerbose }
            ),
            React.createElement(
                Table,
                { key: 'table', groups: current.Groups },
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
