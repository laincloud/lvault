import Cookie from './Cookie';

export let User = {
  
  getSecrets(formData, callback) {
	const {access, ty} = Admin.getTokenCookie()
    const token = access
    const tokenType = ty
    if(formData.procname == undefined){
      Fetch.text(`/v2/secrets?app=${formData.appname}`, 'GET',{token, tokenType} , null, (code, txt) => {
		  callback && callback(code === 200, code === 200 ? `Succeed! ${txt}` : `Failed! ${txt}`);
        }, (msg) => {
        callback && callback(false, msg); 
      });
    }
    else{
      Fetch.text(`/v2/secrets?app=${formData.appname}&proc=${formData.procname}`, 'GET',{token, tokenType} , null, (code, txt) => {
		  callback && callback(code === 200, code === 200 ? `Succeed! ${txt}` : `Failed! ${txt}`);
        }, (msg) => {
        callback && callback(false, msg); 
      });
    }
  },

  putSecret(formData, callback) {
	  const {access, ty} = Admin.getTokenCookie()
      const token = access
      const tokenType = ty
	  Fetch.text(`/v2/secrets?app=${formData.appname}&proc=${formData.procname}&path=${formData.fpath}`, 'PUT', {token,tokenType}, formData, (code, data) => {
      if (code === 200) {
        callback && callback(true, `Succeed! ${data.message || data}`);
      } else {
        callback && callback(false, `Failed! ${data || 'Server got hacked'}`);
      }
    }, (msg) => {
      callback && callback(false, msg); 
    });
  },

};

export const kCookieToken = 'SSO_Site_Access';

export let Admin = {
  clientId: 1,
  realm: 'SSO-Site',
  secret: 'admin',

   deleteSecret(token, tokenType, formData, callback) {
	   Fetch.text(`/v2/secrets?app=${formData.appname}&proc=${formData.procname}&path=${formData.fpath}`, 'DELETE', {token,tokenType}, null, (code, data) => {
      callback && callback(code === 200, code === 200 ? "Successfully deleted the secret" : txt);
    }, (msg) => {
      callback && callback(false, msg); 
    });
  },

  listVaultStatus(token, tokenType, callback) {
	  //console.log({token,tokenType})  
	Fetch.json('/v2/vaultstatus', 'GET', {token, tokenType}, null, (code, data) => {
      callback && callback(code === 200 ? data : []);
    }, (msg) => {
      callback && callback([]); 
    });
  },
  
  listLvaultStatus(token, tokenType, callback) {
	  //console.log({token,tokenType})  
	  Fetch.json('/v2/status', 'GET', {token, tokenType}, null, (code, data) => {
      callback && callback(code === 200 ? data : []);
    }, (msg) => {
      callback && callback([]); 
    });
  },


  setTokenCookie(token, tokenType, expires) {
    const tokenValue = {
      access: token,
      ty: tokenType,
    };
    Cookie.set(kCookieToken, tokenValue, {
      path: '/',
    });
  },

  getTokenCookie() {
    const tokens = Cookie.get(kCookieToken);
    if (tokens) {
      try {
        return JSON.parse(tokens);
      } catch (ex) {}  
    }
    return { access: '', ty: '' };
  },

  getToken(authCode, adminArea, checkGroup, callback) {
    let query = { area: adminArea };
    if (checkGroup) {
      query['ag'] = checkGroup;
    }
    const redirectUrl = `${window.location.protocol}//${window.location.host}/v2/spa/admin/authorize?${this.toQuery(query)}`;
    let formData = {
      client_id: this.clientId,
      client_secret: this.secret,
      code: authCode,
      grant_type: 'authorization_code',
      redirect_uri: redirectUrl,
    };
	Fetch.json(`/v2/oauth2/token?${this.toQuery(formData)}`, 'GET', null, null, (code, data) => {
      if (code === 200) {
        this.setTokenCookie(data.access_token, data.token_type, data.expires_in);
        callback && callback(true, data.access_token, data.token_type);
      } else {
        callback && callback(false);
      }
    }, (msg) => {
      callback && callback(false); 
    });
  },

  redirectOauth(adminArea, checkGroup) {
    let query = { area: adminArea };
    if (checkGroup) {
      query['ag'] = checkGroup;
    }
	const redirectUrl = `${window.location.protocol}//${window.location.host}/v2/spa/admin/authorize?${this.toQuery(query)}`;
    let params = {
      response_type: 'code',
      redirect_uri: redirectUrl,
      realm: this.realm,
      client_id: this.clientId,
      scope: 'write:app read:app read:user write:user write:group read:group',
      state: Math.random(),
    };
	const oauthUrl = `${window.location.protocol}//${window.location.host}/v2/oauth2/auth?${this.toQuery(params)}`;
    window.location.href = oauthUrl;
  },

  toQuery(kv) {
    let params = [];
    _.forOwn(kv, (value, key) => {
      params.push(`${key}=${encodeURIComponent(value)}`);
    });
    return params.join("&");
  },

};

let Fetch = {
  
  text(api, method, auth, payload, succCb, errCb) {
    let code = 200;
    let headers = {
      'Accept': 'application/json',
      'Content-Type': 'application/json'
    };
    if (auth && auth.token && auth.tokenType) {
      const {token, tokenType} = auth;
      headers['Authorization'] = `${_.capitalize(tokenType)} ${token}`;
	  headers['access-token'] = `${token}`
	}
    let options = { method, headers };
    if (payload) {
      options['body'] = JSON.stringify(payload);
    }
    fetch(api, options).then(response => {
      code = response.status;
      return response.text();
    }).then(txt => {
      succCb && succCb(code, txt); 
    }).catch(err => {
      console.log(`Error when ${method} ${api} ${payload}: ${err}`);
      errCb && errCb(`Server got hacked, ${err}`)
    });
  },

  json(api, method, auth, payload, succCb, errCb) {
    this.text(api, method, auth, payload, (code, txt) => {
      succCb && succCb(code, JSON.parse(txt));
    }, (msg) => {
      errCb && errCb(msg);
    })
  },

};

