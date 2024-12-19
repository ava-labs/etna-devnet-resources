import { useEffect, useState } from 'react';
import TimeAgo from 'javascript-time-ago'
import en from 'javascript-time-ago/locale/en'
TimeAgo.addDefaultLocale(en)
const timeAgo = new TimeAgo('en-US')


export const TechInfo = () => {
    const [address, setAddress] = useState<string | null>(null);
    const [pBalance, setPBalance] = useState<string | null>(null);
    const [cBalance, setCBalance] = useState<string | null>(null);
    const [error, setError] = useState<string | null>(null);
    const [isLoading, setIsLoading] = useState(true);
    const [compileTs, setCompileTs] = useState<string | null>(null);

    useEffect(() => {
        const fetchData = async () => {
            try {
                setIsLoading(true);
                const [addrResponse, pBalResponse, cBalResponse, compiledResponse] = await Promise.all([
                    fetch('/api/addr/c'),
                    fetch('/api/balance/p'),
                    fetch('/api/balance/c'),
                    fetch('/api/compiled')
                ]);

                if (!addrResponse.ok || !pBalResponse.ok || !cBalResponse.ok || !compiledResponse.ok) {
                    throw new Error('Failed to fetch data');
                }

                const [addr, pBal, cBal, compiled] = await Promise.all([
                    addrResponse.text(),
                    pBalResponse.text(),
                    cBalResponse.text(),
                    compiledResponse.text()
                ]);

                setAddress(addr);
                setPBalance(parseFloat(pBal).toFixed(2));
                setCBalance(parseFloat(cBal).toFixed(2));
                setCompileTs(compiled);
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
            {address}<br />
            P-Chain: {pBalance} AVAX, C-Chain {cBalance} AVAX
            <br />
            Current release compiled {compileTs ? timeAgo.format(new Date(parseInt(compileTs) * 1000)) : ''}
        </div>
    );
};
