import { NextRequest, NextResponse } from 'next/server';

type Test = {
    message: string;
    error: Error;
};

export const GET = async (req: NextRequest) => {
    return NextResponse.json({ message: 'test' }, { status: 200 });
};
