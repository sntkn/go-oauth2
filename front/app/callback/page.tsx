'use client'
import { useSearchParams } from 'next/navigation'
import { useEffect } from 'react'

export default function Callback() {
  const searchParams = useSearchParams()

  useEffect(() => {
    // クエリパラメータから認可コードを取得
    const code = searchParams.get('code')

    if (code) {
      // 認可コードが取得できた場合、アクセストークンの取得リクエストを送信
      fetch('/api/fetchToken', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ code }),
      })
        .then((response) => response.json())
        .then((data) => {
          // アクセストークンの処理
          console.log('Received access token:', data.access_token)
        })
        .catch((error) => {
          console.error('Error fetching access token:', error)
        })
    } else {
      console.error('Authorization code not found in query parameters.')
    }
  }, [])

  return (
    <div>
      <p>Processing...</p>
    </div>
  )
}
