import { getAuthorizedAuth } from '$lib/firebase/firebaseAdmin';
import type { Metadata } from 'next';
import CustomEventListener from '../event/[id]/components/helpers/CustomEventListener';
import Favorites from './components/Favorites';

type Props = {};

export const metadata: Metadata = {
    title: 'Interesser',
};

const page = async ({}: Props) => {
    const { user } = await getAuthorizedAuth();
    if (!user) {
        return <>Du må være logget inn for å se dine deltakers interesser.</>;
    }
    return (
        <>
            <Favorites />
            <CustomEventListener />
        </>
    );
};

export default page;
