import { buildSchema, graphql } from 'graphql';
import { type NextRequest, NextResponse } from 'next/server';
import { adminDb, adminUser } from '@/lib/firebaseAdmin';
import { FirebaseCollections } from '@/models/enums';
import { CrmJson } from '@/models/types';

// const url = https://app.checkin.no/graphql?client_id=API_KEY&client_secret=API_SECRET
const query = `{
  allCrms(customerId: 13446) {
    records
    data {
      id
      firstName
      lastName
      hash
      email {
        email
        id
      }
    }
  }
}`;
export const GET = async (request: NextRequest, { params: { id } }: { params: { id: string } }) => {
    try {
        await adminUser.getUserByEmail('post:email');
    } catch (error) {
        const res = await fetch(
            `https://app.checkin.no/graphql?client_id=${process.env.CHECKIN_KEY}&client_secret=${process.env.CHECKIN_SECRET}`,
            {
                method: 'POST',
                body: JSON.stringify({ query }),
                headers: {
                    'Content-Type': 'application/json',
                },
            }
        );
        const queryResult: CrmJson | undefined = await res.json();
        if (queryResult?.errors) {
            return NextResponse.json({ errors: queryResult?.errors }, { status: 403 });
        }
        const paidUser = queryResult?.data.allCrms.data.find((crm) => crm?.email?.email === 'regncon@gmail.com');
        console.log(paidUser);

        if (paidUser) {
            const userRecord = await adminUser.createUser({
                email: 'regncon@gmail.com',
                password: 'regncon@gmail.com',
                displayName: `${paidUser.firstName} ${paidUser.lastName}`,
            });
            const customToken = await adminUser.createCustomToken(userRecord.uid);
            console.log(customToken);

            await adminDb
                .collection(FirebaseCollections.userSetting)
                .doc(userRecord.uid)
                .set({ name: `${paidUser.firstName} ${paidUser.lastName}`.trim() });
            return NextResponse.json({ user: userRecord, status: 'created', customToken }, { status: 200 });
        }
        return NextResponse.json({ user: '', status: `Don't exist`, customToken: '' }, { status: 200 });
    }
    return NextResponse.json({ user: '', status: 'Exists', customToken: '' }, { status: 200 });
};

export const POST = async (request: Request) => {
    const res = await request.json();
    return NextResponse.json({ res });
};
