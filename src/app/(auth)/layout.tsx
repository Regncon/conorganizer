import { getAuthorizedAuth } from '$lib/firebase/firebaseAdmin';
import { redirect } from 'next/navigation';

type Props = {
    children: React.ReactNode;
};
const layout = async ({ children }: Props) => {
    const { auth } = await getAuthorizedAuth();
    if (auth === null) {
        return <>{children}</>;
    }
    redirect('/dashboard');
};

export default layout;
