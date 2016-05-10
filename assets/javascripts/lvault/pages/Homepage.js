import StyleSheet from 'react-style';
import React from 'react';

import WriteSecretCard from '../components/WriteSecretCard';
import ReadSecretsCard from '../components/ReadSecretsCard';
import AdminIndexCard from '../components/AdminIndexCard';

let Homepage = React.createClass({

  render() {
    return (
      <div className="mdl-grid">
        <div className="mdl-cell mdl-cell--4-col mdl-cell--8-col-tablet mdl-cell--4-col-phone">
			<WriteSecretCard />
        </div>

        <div className="mdl-cell mdl-cell--4-col mdl-cell--8-col-tablet mdl-cell--4-col-phone">
			<ReadSecretsCard />
        </div>

        <div className="mdl-cell mdl-cell--4-col mdl-cell--8-col-tablet mdl-cell--4-col-phone">
          <AdminIndexCard />
        </div>
      </div>
    );
  },

});

export default Homepage;
