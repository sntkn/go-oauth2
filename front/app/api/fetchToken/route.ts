import { NextResponse, NextRequest } from 'next/server';

type Token = {
  accessToken: string;
  refreshToken: string;
  expiry: number;
}

export async function POST(req: NextRequest) {
  const data = await req.json();
  const code = data.code;
  const res = await fetch('http://localhost:8080/token', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      grant_type: 'authorization_code',
      code,
    })
  })
  const token: Token = await res.json()
  console.log(token)

  // TODO: setcookie

  return NextResponse.json({
    result: true
  })
}
