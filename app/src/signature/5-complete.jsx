import * as React from "react"

const Complete = ({ username, name }) => {
  return (
    <div className="signature-complete">
      <div className="signature-complete-text">
        <strong>{username}</strong>, you've succesfully uploaded a signature named{' '}
        <em>{name}</em>.<br /> Please close this window and return to the Keybase
        client to continue using notarybot.
      </div>
    </div>
  )
}

export default Complete
