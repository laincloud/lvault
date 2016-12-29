import StyleSheet from 'react-style';
import React from 'react';
import {History} from 'react-router';

import {User} from '../models/Models';
import CardFormMixin from './CardFormMixin';

let ReadSecretsCard = React.createClass({
  mixins: [CardFormMixin, History],

  getInitialState() {
    return {
      formValids: {
		  'appname': true,
      },
      appname: "",
      procname: "",
    };
  },

  componentDidUpdate() {
    componentHandler.upgradeDom();
  },

  render() {
    const {reqResult} = this.state;
    return (
      <div className="mdl-card mdl-shadow--2dp" styles={[this.styles.card, this.props.style]}>
        <div className="mdl-card__title">
		<h2 className="mdl-card__title-text">查询秘密文件</h2>
        </div>
        { this.renderResult() }
        {
          reqResult.fin && reqResult.ok ? null :
            this.renderForm(this.onReset, [
				this.renderInput("appname", "APP name*(字母、数字、减号和'.')", { type: "text", pattern: "[\-a-zA-Z0-9.]*" }),
				this.renderInput("procname", "proc 全名(eg appname.web.procname)", { type: 'text' }),
            ])
        }
        { this.renderAction("查询", this.onReset) }
      </div>
    );
  },

  onReset() {
    const {isValid, formData} = this.validateForm(["appname","procname"], ["appname"]);
    if (isValid) {
      this.setState(
        { inRequest: true,
          appname: formData['appname'],
          procname: formData['procname'],
      });
      User.getSecrets(formData, this.onRequestCallbackToAnotherPage);
    }
  },

  onRequestCallbackToAnotherPage(ok, status){
    if (!ok){
      this.setState({
        inRequest: false,
        reqResult: { fin: true, ok, status },
      });
    }
    else{
      console.log("hello")
      this.history.pushState({
        status,                
      }, `/v2/spa/secret/${this.state.appname}/detail/${this.state.procname}`);
    }
  },

});

export default ReadSecretsCard;
