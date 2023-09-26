import { type NextRequest, NextResponse } from 'next/server';
import { adminDb, adminDoc } from '@/lib/firebaseAdmin';

type Test = {
    message: string;
    error: Error;
};

export const GET = async (request: NextRequest, { params: { id } }: { params: { id: string } }) => {
    console.log(id);

    // query is "hello" for /api/search?query=hello
    const searchParams = request.nextUrl.searchParams;
    const dbData = await adminDb.collection('usersettings').get();
    console.log(dbData);

    return NextResponse.json({ message: dbData }, { status: 200 });
};

export const POST = async (request: Request) => {
    const res = await request.json();
    return NextResponse.json({ res });
};
