import * as React from 'react'

import Drop from './1-drop'
import Position from './2-position'

import './style.css'

const Document = ({ username, id, hash }) => {
  // drop
  const [mode, setMode] = React.useState('drop')

  // drop vars
  const [document, setDocument] = React.useState('')
  const onDrop = React.useCallback(data => {
    setDocument(data)
    setMode('position')
  }, [setDocument])

  return (
    <div className="document-wrapper">
      <div className="document-modal">
        <div className="document-header">
          {
            mode === 'drop' ? 'Upload a new document' :
            mode === 'position' ? 'Edit signature fields' :
            'Invalid mode'
          }
        </div>
        <div className="document-body">
          {mode === 'drop' && <Drop onDrop={onDrop} />}
          {mode === 'position' && <Position document={document} />}
        </div>
      </div>
    </div>
  )
}

export default Document
