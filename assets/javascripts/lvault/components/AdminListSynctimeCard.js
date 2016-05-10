import StyleSheet from 'react-style';
import React from 'react';
import {History} from 'react-router';

import {Admin} from '../models/Models';

let AdminListSynctimeCard = React.createClass({
  mixins: [History],

  getInitialState() {
    return {
      syncs: [], 
    };
  },

  componentDidMount() {
    this.reload(); 
  },

  render() {
    return (
      <div className="mdl-card mdl-shadow--2dp" styles={[this.styles.card, this.props.style]}>
        <div className="mdl-card__title">
			<h2 className="mdl-card__title-text">最新 sync_secret 写数据时间列表</h2>
        </div>

        <table className="mdl-data-table mdl-js-data-table" style={this.styles.table}>
          <thead>
            <tr>
				<th className="mdl-data-table__cell--non-numeric">IP</th>
              <th className="mdl-data-table__cell--non-numeric">最近写入时间</th>
            </tr>
          </thead>
          <tbody>
            {
              this.state.syncs.map((group, index) => {
                return (
                  <tr key={`group-${index}`}>
                     <td className="mdl-data-table__cell--non-numeric">{group.ip}</td>
                     <td className="mdl-data-table__cell--non-numeric">{group.time}</td>
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
    Admin.listSyncStatus(token, tokenType, (syncs) => {
      this.setState({ syncs });
    });  
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

export default AdminListSynctimeCard;
