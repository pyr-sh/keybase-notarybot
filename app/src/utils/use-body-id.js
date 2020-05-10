import {useEffect, useState} from 'react'

export default function useBodyID(id) {
  const [originalBodyID, setOriginalBodyID] = useState()

  useEffect(() => {
    const body = document.getElementsByTagName("body")[0]
    setOriginalBodyID(body.id)
    body.id = id
  }, [id, setOriginalBodyID])

  useEffect(() => () => {
    if (originalBodyID) {
      document.getElementsByTagName("body")[0].id = originalBodyID
    }
  }, [originalBodyID])
}
