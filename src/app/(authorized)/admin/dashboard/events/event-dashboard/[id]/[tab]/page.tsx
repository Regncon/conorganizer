import { notFound } from 'next/navigation';
import Edit from './edit/Edit';
import type { TabNames } from './lib/types/types';
import Room from './room/Room';
import Settings from './settings/Settings';

type Props = {
    params: {
        id: string;
        tab: TabNames;
    };
};

const page = ({ params: { id, tab } }: Props) => {
    if (tab === 'edit') {
        return <Edit id={id} />;
    }
    if (tab === 'room') {
        return <Room id={id} />;
    }
    if (tab === 'settings') {
        return <Settings id={id} />;
    }
    return notFound();
};

export default page;
