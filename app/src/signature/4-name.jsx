import * as React from "react"

const Name = ({ value, onChange }) => {
  const handleChange = React.useCallback(
    (e) => onChange(e.target.value.replace(/\W/g, '')),
    [onChange]
  )

  return (
    <div className="signature-name">
      <label htmlFor="signature-name-input" className="signature-name-label">
        <span>Please name your signature:</span>
        <input type="text" value={value} onChange={handleChange} id="signature-name-input" />
      </label>
    </div>
  )
}

export default Name
