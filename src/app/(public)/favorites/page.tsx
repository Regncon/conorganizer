import { getAuthorizedAuth } from '$lib/firebase/firebaseAdmin';
import Favorites from './components/Favorites';

type Props = {};

const page = async ({}: Props) => {
    const { user } = await getAuthorizedAuth();
    if (!user) {
        return <>Du må være logget inn for å se dine favoritter</>;
    }
    return <Favorites />;
};

export default page;
