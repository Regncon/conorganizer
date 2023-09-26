import { NextRequest, NextResponse } from 'next/server';

type Test = {
    message: string;
    error: Error;
};

export async function GET(req: NextRequest) {
    return NextResponse.json({ message: 'test' }, { status: 200 });
}
