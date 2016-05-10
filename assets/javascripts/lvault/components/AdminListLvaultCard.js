import StyleSheet from 'react-style';
import React from 'react';

import {Admin} from '../models/Models';

let AdminListLvaultCard = React.createClass({

  getInitialState() {
    return {
		vaults: [],
		lvaults: [],
    };
  },

  componentDidMount() {
    this.reload(); 
  },

  render() {
    return (
      <div className="mdl-card mdl-shadow--2dp" styles={[this.styles.card, this.props.style]}>
        <div className="mdl-card__title">
			<h2 className="mdl-card__title-text">Lvault 节点状态</h2>
        </div>

        <table className="mdl-data-table mdl-js-data-table" style={this.styles.table}>
          <thead>
            <tr>
				<th className="mdl-data-table__cell--non-numeric">Vault node IP</th>
				<th className="mdl-data-table__cell--non-numeric">端口</th>
			  <th className="mdl-data-table__cell--non-numeric">未解锁</th>
            </tr>
          </thead>
          <tbody>
            {
				this.state.vaults.map((app, index) => {
                return (
                  <tr key={`app-${index}`}>
                    <td className="mdl-data-table__cell--non-numeric">{app.Info.container_ip}</td>
                    <td className="mdl-data-table__cell--non-numeric">{app.Info.container_port}</td>
                    <td className="mdl-data-table__cell--non-numeric">{app.Status.sealed.toString()}</td>
                  </tr>
                );
              })
            }
          </tbody>
        </table>

        <table className="mdl-data-table mdl-js-data-table" style={this.styles.table}>
          <thead>
            <tr>
				<th className="mdl-data-table__cell--non-numeric">Lvault node</th>
				<th className="mdl-data-table__cell--non-numeric">不可用</th>
            </tr>
          </thead>
          <tbody>
            {
				this.state.lvaults.map((app, index) => {
                return (
                  <tr key={`app-${index}`}>
                    <td className="mdl-data-table__cell--non-numeric">{app.Host}</td>
                    <td className="mdl-data-table__cell--non-numeric">{app.IsMiss.toString()}</td>
                  </tr>
                );
              })
            }
          </tbody>
        </table>

        <div className="mdl-card__actions">
          <button className="mdl-button mdl-js-button mdl-js-ripple-effect mdl-button--colored"
            onClick={this.reload}>刷新</button>
        </div>
      </div>
    );
  },

  reload() {
    const {token, tokenType} = this.props;
	console.log("begin reload");
	Admin.listVaultStatus(token, tokenType, (vaults) => {
      this.setState({ vaults });
	});
	Admin.listLvaultStatus(token, tokenType, (lvaults) => {
		this.setState({ lvaults });
	});
	console.log("after reload");
  },

  styles: StyleSheet.create({
    card: {
      width: '100%',
      marginBottom: 16,
      minHeight: 50,
    },
    table: {
      width: '100%',
      borderLeft: 'none',
      borderRight: 'none',
    },
  }),

});

export default AdminListLvaultCard;
