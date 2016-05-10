import StyleSheet from 'react-style';
import React from 'react';

import CardFormMixin from './CardFormMixin';
import {Admin} from '../models/Models';

let AdminDeleteSecretCard = React.createClass({
  mixins: [CardFormMixin],

  getInitialState() {
    return {
      formValids: {
		  'appname': true,
		  'procname': true,
		  'fpath': true
      },
    };
  },

  render() {
    const {reqResult} = this.state;
    return (
      <div className="mdl-card mdl-shadow--2dp" styles={[this.styles.card, this.props.style]}>
        <div className="mdl-card__title">
          <h2 className="mdl-card__title-text">删除秘密文件</h2>
        </div>
        { this.renderResult() }
        { 
          reqResult.fin && reqResult.ok ? null :
            this.renderForm(this.onDelete, [
               this.renderInput("appname", "APP name*(字母、数字、减号和'.')", { type: "text", pattern: "[\-a-zA-Z0-9.]*" }),
              this.renderInput("procname", "proc 全名*(eg appname.web.procname)", { type: 'text' }),
			  this.renderInput("fpath", "绝对路径*(eg /lain/app/hello.dat)", { type: 'text' }),
            ])
        }
        { this.renderAction("确定删除", this.onDelete) }
      </div>
    );
  },

  onDelete() {
    const fields = ['appname','procname','fpath'];
    const {isValid, formData} = this.validateForm(fields, fields);
    if (isValid) {
      const {token, tokenType} = this.props;
      this.setState({ inRequest: true });
      Admin.deleteSecret(token, tokenType, formData, this.onRequestCallback);
    }
  },

});

export default AdminDeleteSecretCard;
