import { type NextRequest, NextResponse } from 'next/server';
import { adminDb } from '@/lib/firebaseAdmin';
import { FirebaseCollections } from '@/models/enums';

type Test = {
    message: string;
    error: Error;
};

export const GET = async (request: NextRequest, { params: { id } }: { params: { id: string } }) => {
    // her må vi få en array og loope igjennom for og oppdatere dokumenter
    // for eksempel bruke adminDb.batch
    // const batch = adminDb.batch();
    // const ids = ['123345123', '12312312345']
    // ids.forEach((id) => {
    //     batch.set(adminDb.doc(`${FirebaseCollections.Test}/${id}`), { test: 'abc123' });
    // })
    // try {
    //     dbData = await batch.commit();
    // } catch (error) {
    //     NextResponse.json({ Error: `Failed to set data in id: ${id}` }, { status: 403 });
    // }
    let dbData;
    try {
        dbData = await adminDb
            .collection(FirebaseCollections.Test)
            .doc(id + Math.floor(Math.random() * 999))
            .get();
    } catch (error) {
        NextResponse.json({ Error: `Failed to set data in id: ${id}` }, { status: 403 });
    }
    console.log(dbData);

    console.log(dbData);

    return NextResponse.json({ message: dbData }, { status: 200 });
};

export const POST = async (request: Request) => {
    const res = await request.json();
    return NextResponse.json({ res });
};
