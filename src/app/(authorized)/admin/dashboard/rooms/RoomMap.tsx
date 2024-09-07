import { Box } from '@mui/material';
import Image from 'next/image';
import RoomCard from './RoomCard';
import { PoolName, RoomName } from '$lib/enums';

type Props = {
    pool: PoolName;
};

const RoomMap = async ({ pool }: Props) => {
    return (
        <Box>
            <RoomCard
                poolName={pool}
                roomName={RoomName.Klang}
                title={'Kjempegøy drager og fangehull'}
                gameMaster={'Kari Nordmann'}
                system={'D&D'}
                imageUri="/blekksprut2.jpg"
            ></RoomCard>
            <RoomCard
                poolName={pool}
                roomName={RoomName.Sonate}
                title={'En telefon fra Cthulhu'}
                gameMaster={'Ola Nordmann'}
                system={'Call of Cthulhu'}
                imageUri="/blekksprut2.jpg"
            ></RoomCard>
            <Image src={'/rooms.webp'} alt={'Romkart'} width={'2901'} height={'2073'}></Image>
        </Box>
    );
};
export default RoomMap;
