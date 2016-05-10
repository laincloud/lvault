import StyleSheet from 'react-style';
import React from 'react';
import {History} from 'react-router';

import {Admin} from '../models/Models';

let ListSecretsCard = React.createClass({
  mixins: [History],

  render() {
    return (
      <div className="mdl-card mdl-shadow--2dp" styles={[this.styles.card, this.props.style]}>
        <div className="mdl-card__title">
			<h2 className="mdl-card__title-text">secret files 列表</h2>
        </div>

        <table className="mdl-data-table mdl-js-data-table" style={this.styles.table}>
          <thead>
            <tr>
				<th className="mdl-data-table__cell--non-numeric">路径</th>
              <th className="mdl-data-table__cell--non-numeric">内容</th>
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

      </div>
    );
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

export default ListSecretsCard;
