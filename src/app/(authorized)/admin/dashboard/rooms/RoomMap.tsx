import { Box, Typography } from '@mui/material';
import Image from 'next/image';
import RoomCard from './RoomCard';
import { PoolName, RoomName } from '$lib/enums';

type Props = {
    pool: PoolName;
};

const RoomMap = async ({ pool }: Props) => {
    let poolTitle = '';
    if (pool === PoolName.fridayEvening) {
        poolTitle = 'Fredag Kveld';
    }
    if (pool === PoolName.saturdayMorning) {
        poolTitle = 'Lørdag Morgen';
    }
    if (pool === PoolName.saturdayEvening) {
        poolTitle = 'Lørdag Kveld';
    }
    if (pool === PoolName.sundayMorning) {
        poolTitle = 'Søndag Morgen';
    }

    return (
        <Box>
            <Typography
                variant="h1"
                sx={{ fontSize: '90px', color: 'black', position: 'absolute', top: '100px', left: '900px' }}
            >
                {poolTitle}
            </Typography>
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
