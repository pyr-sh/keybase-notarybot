import * as React from "react";

const Complete = ({ username, id }) => {
  return (
    <div className="document-complete">
      <div className="document-complete-text">
        <strong>{username}</strong>, you've succesfully uploaded a document
        named <em>{id}</em>.<br /> Please close this window and return to the
        Keybase client to continue using notarybot.
      </div>
    </div>
  );
};

export default Complete;
