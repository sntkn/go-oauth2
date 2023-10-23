import { NextResponse } from 'next/server';

export async function POST(request: Request) {
  return NextResponse.json({
    access_token: 'testaaaa',
    refresh_token: 'test_refresh_token',
    expiry: 12345688
  });
}
