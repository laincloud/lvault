import 'es5-shim';
import 'whatwg-fetch';
import _ from 'lodash';
window._ = _;

import StyleSheet from 'react-style';
import React from 'react';
import {Router, Route} from 'react-router';
import createBrowserHistory from 'history/lib/createBrowserHistory'

import {createElement} from './lvault/app';
import Homepage from './lvault/pages/Homepage';
import WriteSecretPage from './lvault/pages/WriteSecretPage';
import LogoutPage from './lvault/pages/LogoutPage';
import AdminAuthorizePage from './lvault/pages/AdminAuthorizePage';
import AdminLvaultsPage from './lvault/pages/AdminLvaultsPage';
import AdminDeletePage from './lvault/pages/AdminDeletePage';
import ReadSecretsPage from './lvault/pages/ReadSecretsPage';
import ListSecretsPage from './lvault/pages/ListSecretsPage';
import EditSecretPage from './lvault/pages/EditSecretPage';

let domReady = () => {
  React.initializeTouchEvents(true);

  let history = createBrowserHistory();
  React.render((
    <Router history={history} >
		<Route path="/v2/spa" component={createElement(Homepage)} />

		<Route path="/v2/spa/secret/write" component={createElement(WriteSecretPage)} />
        <Route path="/v2/spa/secret/read" component={createElement(ReadSecretsPage)} />
        <Route path="/v2/spa/secret/:app/detail/:proc" component={createElement(ListSecretsPage)} />
        <Route path="/v2/spa/secret-edit/:app" component={createElement(EditSecretPage)} />
		<Route path="/v2/spa/user/logout" component={createElement(LogoutPage)} />

		<Route path="/v2/spa/admin/authorize" component={createElement(AdminAuthorizePage)} />
		<Route path="/v2/spa/admin/lvaults" component={createElement(AdminLvaultsPage)} />
		<Route path="/v2/spa/admin/delete" component={createElement(AdminDeletePage)} />
    </Router>
  ), document.getElementById("lvault-spa"));
};

if (typeof document.onreadystatechange === "undefined") {
    window.onload = () => domReady();
} else {
    document.onreadystatechange = () => {
      if (document.readyState !== "complete") {
        return;
      }
      domReady();
    }
}
