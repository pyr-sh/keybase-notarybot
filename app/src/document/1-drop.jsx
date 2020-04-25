import * as React from 'react'
import { useDropzone } from 'react-dropzone'

const Drop = ({ onDrop }) => {
  const handleDrop = React.useCallback(acceptedFiles => {
    if (acceptedFiles.length < 1) {
      return
    }

    const reader = new FileReader()
    reader.onload = e => onDrop(e.target.result)
    reader.readAsDataURL(acceptedFiles[0])
  }, [onDrop])
  const {getRootProps, getInputProps, isDragActive} = useDropzone({
    accept: ['application/pdf'],
    multiple: false,
    onDrop: handleDrop,
  })

  return (
    <div {...getRootProps()} className="document-drop">
      <input {...getInputProps()} />
      {
        isDragActive ?
          <p>Drop the file here...</p> :
          <p>Drag 'n' drop your document here</p>
      }
    </div>
  )
}

export default Drop
