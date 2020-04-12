import * as React from 'react'
import ReactCrop from 'react-image-crop'

import 'react-image-crop/dist/ReactCrop.css'

export const getCroppedImage = (image, crop) => {
  return new Promise(resolve => {
    const img = new Image()
    img.onload = () => {
      const canvas = document.createElement('canvas')
      const croppedWidth = crop.width / 100 * img.naturalWidth
      const croppedHeight = crop.height / 100 * img.naturalHeight
      canvas.width = croppedWidth
      canvas.height = croppedHeight
      const ctx = canvas.getContext('2d')

      ctx.drawImage(
        img,
        crop.x / 100 * img.naturalWidth,
        crop.y / 100 * img.naturalHeight,
        croppedWidth,
        croppedHeight,
        0,
        0,
        croppedWidth,
        croppedHeight,
      )
      return resolve(canvas.toDataURL('image/png'))
    }
    img.src = image
  })
}

const Crop = ({image, crop, onCrop}) => {
  return (
    <ReactCrop src={image} crop={crop} onChange={(_, newCrop) => onCrop(newCrop)} />
  )
}

export default Crop
