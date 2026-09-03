import { useEffect, useRef } from 'react'
import QrScannerLib from 'qr-scanner'

interface QrScannerProps {
  onScan: (payload: string) => void
  onError: (message: string) => void
}

export function QrScanner({ onScan, onError }: QrScannerProps) {
  const videoRef = useRef<HTMLVideoElement>(null)
  const scannerRef = useRef<QrScannerLib | null>(null)

  useEffect(() => {
    const video = videoRef.current
    if (!video) {
      return
    }

    const scanner = new QrScannerLib(
      video,
      (result) => {
        scanner.stop()
        onScan(result.data)
      },
      { returnDetailedScanResult: true },
    )
    scannerRef.current = scanner

    scanner.start().catch(() => {
      onError('Could not access the camera. Check permissions and try again.')
    })

    return () => {
      scanner.destroy()
      scannerRef.current = null
    }
  }, [onScan, onError])

  return (
    <div role="img" aria-label="Camera view for scanning a client QR code">
      <video ref={videoRef} className="w-full rounded-card" />
    </div>
  )
}
