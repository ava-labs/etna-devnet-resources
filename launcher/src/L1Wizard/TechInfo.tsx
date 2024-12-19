import { useEffect, useState } from 'react';

export const TechInfo = () => {
    const [address, setAddress] = useState<string | null>(null);
    const [pBalance, setPBalance] = useState<string | null>(null);
    const [cBalance, setCBalance] = useState<string | null>(null);
    const [error, setError] = useState<string | null>(null);
    const [isLoading, setIsLoading] = useState(true);

    useEffect(() => {
        const fetchData = async () => {
            try {
                setIsLoading(true);
                const [addrResponse, pBalResponse, cBalResponse] = await Promise.all([
                    fetch('/api/addr/c'),
                    fetch('/api/balance/p'),
                    fetch('/api/balance/c')
                ]);

                if (!addrResponse.ok || !pBalResponse.ok || !cBalResponse.ok) {
                    throw new Error('Failed to fetch data');
                }

                const [addr, pBal, cBal] = await Promise.all([
                    addrResponse.text(),
                    pBalResponse.text(),
                    cBalResponse.text()
                ]);

                setAddress(addr);
                setPBalance(pBal);
                setCBalance(cBal);
                setError(null);
            } catch (err) {
                setError('Failed to load faucet data');
            } finally {
                setIsLoading(false);
            }
        };

        fetchData();
    }, []);

    if (isLoading) {
        return <div>Loading faucet information...</div>;
    }

    if (error) {
        return <div>Error: {error}</div>;
    }

    return (
        <div>
            Faucet C-Chain address: {address}<br />
            Balances: P-Chain {pBalance} AVAX, C-Chain {cBalance} AVAX
        </div>
    );
};
