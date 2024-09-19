import { getAuthorizedAuth } from '$lib/firebase/firebaseAdmin';
import type { IdTokenResult, ParsedToken } from 'firebase/auth';
import MainAppBarDesktop from './ui/mainAppBarDesktop';
import MainAppBarMobile from './ui/mainAppBarMobile';

type ParsedTokenWithAdmin = (ParsedToken & { admin: boolean | undefined }) | undefined;
type Props = {};

const MainAppBar = async ({}: Props) => {
    const { user } = await getAuthorizedAuth();
    const claims = (await user?.getIdTokenResult())?.claims as ParsedTokenWithAdmin;
    const admin = claims?.admin ?? false;

    return (
        <>
            <MainAppBarDesktop isLoggedIn={!!user} admin={admin} />
            <MainAppBarMobile isLoggedIn={!!user} admin={admin} />
        </>
    );
};

export default MainAppBar;
