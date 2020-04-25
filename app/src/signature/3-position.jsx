import * as React from 'react'
import {Rnd} from 'react-rnd'

const Position = ({image, maxWidth, maxHeight, coords, setCoords, size, setSize}) => {
  const handleDragStop = React.useCallback((e, d) => {
    let newX = d.x
    let newY = d.y
    
    if (newX < 0) {
      newX = 0
    }
    if (newY < 0) {
      newY = 0
    }

    if (newX + size[0] > maxWidth) {
      newX = maxWidth - size[0]
    }
    if (newY + size[1] > maxHeight) {
      newY = maxHeight - size[1]
    }

    setCoords([newX, newY])
  }, [size, maxWidth, maxHeight, setCoords])

  const handleResizeStop = React.useCallback((e, direction, ref, delta, position) => {
    let newX = position.x
    let newY = position.y

    let newWidth = size[0] + delta.width
    let newHeight = size[1] + delta.height

    if (newX + newWidth > maxWidth) {
      newWidth = maxWidth - newX
    }
    if (newY + newHeight > maxHeight) {
      newHeight = maxHeight - newY
    }

    setCoords([newX, newY])
    setSize([newWidth, newHeight])
  }, [size, maxWidth, maxHeight, setSize, setCoords])

  React.useEffect(() => {
    const img = document.createElement('img')
    img.onload = () => {

      const fieldRatio = maxWidth / maxHeight
      const imageRatio = img.width / img.height

      let width = 0
      let height = 0

      if (img.width > maxWidth || img.height > maxHeight) {
        if (fieldRatio === imageRatio) {
          width = maxWidth
          height = maxHeight
        } else if (imageRatio > fieldRatio) {
          width = maxWidth
          height = maxWidth / imageRatio
        } else {
          width = maxHeight * imageRatio
          height = maxHeight
        }
      } else {
        width = maxWidth
        height = maxHeight
      }

      if (maxWidth - width < 30) {
        width = maxWidth - 30
        height = width / imageRatio
      }

      if (maxHeight - height < 30) {
        height = maxHeight - 30
        width = height * imageRatio
      }

      setCoords([maxWidth / 2 - width / 2, maxHeight / 2 - height / 2])
      setSize([width, height])
    }
    img.src = image
  }, [image, setCoords, setSize, maxHeight, maxWidth])

  if (size[0] === 0 || size[1] === 0) {
    return (
      <div className="signature-preview">
        <span className="signature-preview-loading">
          Loading...
        </span>
      </div>
    )
  }

  return (
    <div className="signature-preview">
      <Rnd
        className="signature-preview-draggable"
        position={{x: coords[0], y: coords[1]}}
        size={{width: size[0], height: size[1]}}
        onDragStop={handleDragStop}
        onResizeStop={handleResizeStop}
        lockAspectRatio={true}
      >
        <img src={image} alt="Signature" />
      </Rnd>

      <div className="signature-preview-overlay">
        <svg height="4" width="600">
          {[...Array(59).keys()].map((i) => <circle cx={(i+1)*10} cy={2} r={2} key={i} />)}
        </svg>
      </div>
    </div>
  )
}

export default Position
