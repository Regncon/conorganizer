'use client';
import Card from '@mui/material/Card';
import CardActionArea, { cardActionAreaClasses } from '@mui/material/CardActionArea';
import CardMedia from '@mui/material/CardMedia';
import CardContent from '@mui/material/CardContent';
import Typography from '@mui/material/Typography';
import { useRouter } from 'next/navigation';
import { useEffect, useState } from 'react';
import type { Route } from 'next';
import { iconButtonClasses } from '@mui/material';
type Props = {
    href: Route;
    title: string;
    subTitle: string;
    img: string;
    imgAlt: string;
};

const CardBase = ({ title, img, imgAlt, subTitle, href }: Props) => {
    const router = useRouter();
    const [disableRipple, setDisableRipple] = useState<boolean>(false);
    useEffect(() => {
        router.prefetch(href);
    });
    const handleActionClick = () => {
        router.push(href);
    };

    return (
        <Card sx={{ maxWidth: 345 }}>
            <CardActionArea
                onClick={handleActionClick}
                disableRipple={disableRipple}
                sx={{
                    [`.${cardActionAreaClasses.root}:has(.${iconButtonClasses.root}:hover) .${cardActionAreaClasses.focusHighlight}`]:
                        {
                            backgroundColor: 'unset',
                        },
                }}
            >
                <CardMedia component="img" height={130} image={img} alt={imgAlt} />
                <CardContent>
                    <Typography gutterBottom variant="h5" component="div">
                        {title}
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                        {subTitle}
                    </Typography>
                </CardContent>
            </CardActionArea>
        </Card>
    );
};

export default CardBase;
