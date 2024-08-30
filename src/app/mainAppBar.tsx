import { getAuthorizedAuth } from '$lib/firebase/firebaseAdmin';
import type { IdTokenResult, ParsedToken } from 'firebase/auth';
import MainAppBarDesktop from './mainAppBarDesktop';
import MainAppBarMobile from './mainAppBarMobile';

type ParsedTokenWithAdmin = (ParsedToken & { admin: boolean | undefined }) | undefined;
type Props = {};

const MainAppBar = async ({}: Props) => {
    const { user } = await getAuthorizedAuth();
    const claims = (await user?.getIdTokenResult())?.claims as ParsedTokenWithAdmin;
    const admin = claims?.admin ?? false;

    return (
        <>
            <MainAppBarDesktop user={!!user} admin={admin} />
            <MainAppBarMobile user={!!user} admin={admin} />
        </>
    );
};

export default MainAppBar;
