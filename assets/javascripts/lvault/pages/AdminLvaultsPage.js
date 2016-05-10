import StyleSheet from 'react-style';
import React from 'react';
import {History} from 'react-router';

import AdminListLvaultCard from '../components/AdminListLvaultCard';
import AdminAuthorizeMixin from '../components/AdminAuthorizeMixin';
import {Admin} from '../models/Models';

let AdminLvaultsPage = React.createClass({
  mixins: [History, AdminAuthorizeMixin],

    componentWillMount() {
       this.authorize('lvaults');
	},

  render() {
    const isValid = this.isSessionValid();
    return (
      <div className="mdl-grid">
		  <div className="mdl-cell mdl-cell--12-col mdl-cell--8-col-tablet mdl-cell--4-col-phone">
          {
               <AdminListLvaultCard ref="appList" token={this.state.token} tokenType={this.state.tokenType} />
          }
        </div>
      </div>
    );
  },

  postCreateApp() {
    if (this.isSessionValid()) {
      this.refs.appList.reload();
    }
  },
  
  styles: StyleSheet.create({
  }),

});

export default AdminLvaultsPage;
