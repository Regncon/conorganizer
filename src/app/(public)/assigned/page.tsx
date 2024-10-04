import CustomEventListener from '../event/[id]/components/helpers/CustomEventListener';
import AssignedGame from './components/AssignedGame';

type Props = {};

const page = ({}: Props) => {
    return (
        <>
            <AssignedGame />
            <CustomEventListener />
        </>
    );
};

export default page;
