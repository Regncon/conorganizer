/** @type {import('next').NextConfig} */
const nextConfig = {
    experimental: {
        typedRoutes: true,
    },
    images: {
        remotePatterns: [
            {
                protocol: 'https',
                hostname: 'www.regncon.no',
                port: '',
                pathname: '/regncon2024images/**',
            },
        ],
    },
};

export default nextConfig;
