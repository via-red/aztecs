import React, { useState, useEffect } from 'react'; // Import useEffect
import { Card, Button, Input, Typography, List } from 'antd';
import axios from 'axios'; // Import axios

const { Title } = Typography;

function WalletManager() {
  const [wallets, setWallets] = useState([]); // wallets will now be an array of objects { address: string, balance: number }
  const [newWalletAddress, setNewWalletAddress] = useState('');
  const [loading, setLoading] = useState(false); // Add loading state
  const [fetchingWallets, setFetchingWallets] = useState(true); // Add state for fetching wallets

  const fetchWallets = async () => { // Add fetchWallets function
    try {
      const response = await axios.get('http://localhost:8080/wallets'); // Call get wallets API
      const addresses = response.data.wallets;

      // Fetch balance for each wallet
      const walletsWithBalance = await Promise.all(addresses.map(async (address) => {
        try {
          const balanceResponse = await axios.get(`http://localhost:8080/wallets/${address}/balance`);
          return { address, balance: balanceResponse.data.balance };
        } catch (balanceError) {
          console.error(`Error fetching balance for ${address}:`, balanceError);
          return { address, balance: 'Error' }; // Handle error case
        }
      }));

      setWallets(walletsWithBalance); // Set wallets with their balances
      setFetchingWallets(false);
    } catch (error) {
      console.error('Error fetching wallets:', error);
      setFetchingWallets(false);
    }
  };

  useEffect(() => {
    fetchWallets(); // Fetch wallets on component mount
  }, []);


  const handleCreateWallet = async () => { // Make function async
    setLoading(true);
    try {
      const response = await axios.post('http://localhost:8080/wallets'); // Call create wallet API
      const createdWalletAddress = response.data.address;
      // setWallets([...wallets, createdWalletAddress]); // Add new wallet to list - instead, refetch
      fetchWallets(); // Refresh wallets after creating a new one
      setNewWalletAddress(createdWalletAddress);
    } catch (error) {
      console.error('Error creating wallet:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <Card>
      <Title level={2}>钱包管理</Title>
      <Button type="primary" onClick={handleCreateWallet} loading={loading} style={{ marginBottom: 16 }}> {/* Add loading state to button */}
        创建新钱包
      </Button>
      {newWalletAddress && (
        <p>新创建的钱包地址: <strong>{newWalletAddress}</strong></p>
      )}
      <Title level={3}>我的钱包</Title>
      <List
        bordered
        dataSource={wallets}
        loading={fetchingWallets} // Add loading state to list
        renderItem={(item) => ( // item is now an object { address, balance }
          <List.Item>
            <div>
              <strong>地址:</strong> {item.address} <br />
              <strong>余额:</strong> {item.balance} AZT
            </div>
            {/* TODO: Add functionality to view wallet details */}
          </List.Item>
        )}
      />
      {/* TODO: Add functionality to import wallet */}
    </Card>
  );
}

export default WalletManager;