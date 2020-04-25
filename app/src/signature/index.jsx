import * as React from 'react'
import clsx from 'clsx'
import axios from 'axios'

import {API_URL} from '../constants'

import Drop from './1-drop'
import Crop, {getCroppedImage} from './2-crop'
import Position from './3-position'
import Complete from './4-complete'

import './style.css'

const maxWidth = 600
const maxHeight = 300

const Signature = ({ username, id, hash }) => {
  // drop, crop, position, complete
  const [mode, setMode] = React.useState('drop')

  // Drag and drop handler, manages the transition between the drop and crop modes
  const [uncroppedImage, setUncroppedImage] = React.useState('')
  const onDrop = React.useCallback(data => {
    setUncroppedImage(data)
    setMode('crop')
  }, [setUncroppedImage])

  // Crop  
  const [crop, setCrop] = React.useState({})
  const [croppedImage, setCroppedImage] = React.useState('')

  // Position
  const [coords, setCoords] = React.useState([0, 0])
  const [size, setSize] = React.useState([0, 0])

  // The rest of the flow is handled here and in the handleContinue function
  const canContinue = React.useMemo(() => {
    if (mode === 'crop') {
      // We need a proper crop selection
      if (!crop || isNaN(crop.x) || isNaN(crop.y) || isNaN(crop.width) || isNaN(crop.height) || crop.width === 0 || crop.height === 0) {
        return false
      }
      return true
    }

    if (mode === 'position') {
      return coords[0] !== 0 && coords[1] !== 0
    }

    return false
  }, [mode, crop, coords])
  const handleContinue = React.useCallback(async () => {
    if (mode === 'crop') {
      const image = await getCroppedImage(uncroppedImage, crop)
      setCroppedImage(image)
      setMode('position')
      return
    }

    if (mode === 'position') {
      // The idea at this point is to draw the line over the signature
      const img = new Image()
      img.onload = async () => {
        const canvas = document.createElement('canvas')
        canvas.width = img.naturalWidth
        canvas.height = img.naturalHeight
        const ctx = canvas.getContext('2d')

        // max width is 600, max height is 300
        // line is going through the middle, offset with 50px, so 150+50=200
        ctx.drawImage(
          img,
          0,
          0,
          img.naturalWidth,
          img.naturalHeight,
          0,
          0,
          img.naturalWidth,
          img.naturalHeight,
        )

        let percentageHeight = null

        if (coords[1] + size[1] > 200) {
          // Calculate percentage-wise at what height the dotted line is cutting through the image
          const distanceFromTop = 200 - coords[1]
          percentageHeight = 1 - (size[1] - distanceFromTop) / size[1]
        }

        // Send the request to the API
        try {
          await axios({
            method: 'post',
            url: API_URL + '/signatures',
            data: croppedImage,
            params: {
              "u":    username,
              "id":   id,
              "sig":  hash,
              "p":    percentageHeight,
            },
          })
          setMode('complete')
        } catch (e) {
          alert(e.response.data.error)
        }
      }
      img.src = croppedImage
    }
  }, [mode, uncroppedImage, setCroppedImage, croppedImage, crop, coords, size, hash, id, username])


  return (
    <div className="signature-wrapper">
      <div className="signature-modal">
        <div className="signature-header">
          {
            mode === 'drop' ? 'Upload a signature' :
            mode === 'crop' ? 'Crop the signature' :
            mode === 'position' ? 'Position the signature on the dotted line' :
            mode === 'complete' ? 'Success!' :
            'Invalid mode'
          }
        </div>

        <div className="signature-body">
          {mode === 'drop' && <Drop onDrop={onDrop} />}
          {mode === 'crop' && <Crop image={uncroppedImage} crop={crop} onCrop={setCrop} />}
          {mode === 'position' && <Position
            image={croppedImage}
            maxWidth={maxWidth}
            maxHeight={maxHeight}
            size={size}
            setSize={setSize}
            coords={coords}
            setCoords={setCoords}
          />}
          {mode === 'complete' && <Complete username={username} id={id} />}
        </div>
        <div className="signature-actions">
          <button
            disabled={mode === 'drop' || mode === 'complete'}
            className={clsx('upload-again', {
              disabled: mode === 'upload' || mode === 'complete',
            })}
            onClick={() => setMode('drop')}
          >Upload again</button>
          <button
            disabled={mode === 'drop' || mode === 'crop' || mode === 'complete'}
            className={clsx('upload-again', {
              disabled: mode === 'drop' || mode === 'crop' || mode === 'complete',
            })}
            onClick={() => setMode('crop')}
          >Crop again</button>
          <button
            disabled={!canContinue}
            className={clsx('continue', {
              disabled: !canContinue,
            })}
            onClick={handleContinue}
          >Continue</button>
        </div>
      </div>
    </div>
  )
}

export default Signature
