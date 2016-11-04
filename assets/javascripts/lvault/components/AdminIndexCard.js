import StyleSheet from 'react-style';
import React from 'react';
import {History} from 'react-router';

let AdminIndexCard = React.createClass({
  mixins: [History],
  
  render() {
    const buttons = [
	  { title: "登录 / Lvault 集群状态查询", target: "lvaults" },
      { title: "删除秘密文件", target: "delete" },
    ];
    
    return (
      <div className="mdl-card mdl-shadow--2dp" styles={[this.styles.card, this.props.style]}>
        <div className="mdl-card__title">
          <h2 className="mdl-card__title-text">自助服务</h2>
        </div>
        <div className="mdl-card__supporting-text" style={this.styles.supporting}>
			先注册应用成为应用管理者，然后登录可以写自己应用的 secret files.
		</div>
        
        <div style={{ padding: 8 }}>
          {
            _.map(buttons, (btn) => {
              return (
                <button className="mdl-button mdl-js-button mdl-button--accent mdl-js-ripple-effect"
                  onClick={(evt) => this.adminAuthorize(btn.target)}
                  key={btn.target}
                  style={this.styles.buttons}>
                  {btn.title}
                </button>
              ); 
            })
          }
        </div>
      </div>
    );  
  },

  adminAuthorize(target) {
	  this.history.pushState(null, `/v2/spa/admin/${target}`);
  },

  styles: StyleSheet.create({
    card: {
      width: '100%',
      marginBottom: 16,
      minHeight: 50,
    },
    buttons: {
      display: 'block',
    },
    supporting: {
      borderTop: '1px solid rgba(0, 0, 0, .12)',
      borderBottom: '1px solid rgba(0, 0, 0, .12)',
    },
  }),

});

export default AdminIndexCard;
