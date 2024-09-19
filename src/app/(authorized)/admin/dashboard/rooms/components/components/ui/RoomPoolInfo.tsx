import { PoolName } from '$lib/enums';
import { ConEvent } from '$lib/types';
import { Divider, Stack, Typography } from '@mui/material';
import { poolTitles } from '../../lib/helpers';
import UnwantedPoolByGm from './UnwantedPoolByGm';

type props = {
    poolName: PoolName;
    conEvent: ConEvent;
};

const RoomPoolInfo = ({ poolName, conEvent }: props) => {
    const poolIsInEvent: boolean = conEvent.poolIds?.some((poolId) => poolId.poolName === poolName) ?? false;
    const roomIsInPool: boolean = conEvent.roomIds?.some((roomId) => roomId.poolName === poolName) ?? false;

    if (poolIsInEvent === false) {
        return null;
    }
    return (
        <Stack>
            <Divider />
            <Typography variant="h4">{poolTitles[poolName]} </Typography>
            <UnwantedPoolByGm poolName={poolName} conEvent={conEvent} color="error.main" />
            {roomIsInPool ?
                <>
                    {conEvent.roomIds?.map((roomId) => {
                        if (roomId.poolName === poolName) {
                            return <Typography key={roomId.roomName}>Rom: {roomId.roomName}</Typography>;
                        }
                        return null;
                    })}
                </>
                : <Typography>Ikke tildelt rom</Typography>}
        </Stack>
    );
};

export default RoomPoolInfo;
