'use client';

import MainNavigator from '@/components/Navigation/MainNavigator';
import Participants from './Participants';

type Props = {
    params: { id: string };
};
const page = ({ params }: Props) => {
    const { id } = params;
    return (
        <>
            <Participants />
            <MainNavigator />
        </>
    );
};
export default page;
