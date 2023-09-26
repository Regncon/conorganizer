import { type NextRequest, NextResponse } from 'next/server';

type Test = {
    message: string;
    error: Error;
};

export const GET = async (request: NextRequest, { params: { id } }: { params: { id: number } }) => {
    // query is "hello" for /api/search?query=hello
    const searchParams = request.nextUrl.searchParams;

    if (!id) {
        return NextResponse.json({ message: 'NOOOO NOO U DONT' }, { status: 403 });
    }

    return NextResponse.json({ message: 'test' }, { status: 200 });
};

export const POST = async (request: Request) => {
    const res = await request.json();
    return NextResponse.json({ res });
};
