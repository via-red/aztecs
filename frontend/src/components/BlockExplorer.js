import React, { useEffect, useState } from 'react';
import { Card, Table, Typography, Button } from 'antd'; // Import Button
import axios from 'axios'; // Import axios

const { Title } = Typography;

function BlockExplorer() {
  const [blocks, setBlocks] = useState([]);
  const [loading, setLoading] = useState(true);
  const [mining, setMining] = useState(false); // Add mining state

  const fetchBlockchain = async () => { // Move fetchBlockchain outside useEffect
    try {
      const response = await axios.get('http://localhost:8080/blockchain'); // Call backend API
      // Add a unique key for each block for Ant Design Table
      const blocksWithKeys = response.data.map((block, index) => ({
        ...block,
        key: block.Index, // Use block index as key
        timestamp: new Date(block.Timestamp).toLocaleString(), // Format timestamp
      }));
      setBlocks(blocksWithKeys);
      setLoading(false);
    } catch (error) {
      console.error('Error fetching blockchain:', error);
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchBlockchain();
  }, []);

  const handleMineBlock = async () => { // Add handleMineBlock function
    setMining(true);
    try {
      await axios.post('http://localhost:8080/mine'); // Call mine API
      fetchBlockchain(); // Refresh blockchain data after mining
    } catch (error) {
      console.error('Error mining block:', error);
    } finally {
      setMining(false);
    }
  };

  const columns = [
    { title: 'Index', dataIndex: 'index', key: 'index' },
    { title: 'Timestamp', dataIndex: 'timestamp', key: 'timestamp' },
    { title: 'Data', dataIndex: 'data', key: 'data' },
    { title: 'Previous Hash', dataIndex: 'prevHash', key: 'prevHash' },
    { title: 'Hash', dataIndex: 'hash', key: 'hash' },
    { title: 'Nonce', dataIndex: 'nonce', key: 'nonce' },
  ];

  return (
    <Card>
      <Title level={2}>区块链浏览器</Title>
      <Button type="primary" onClick={handleMineBlock} loading={mining} style={{ marginBottom: 16 }}> {/* Add Mine Button */}
        挖矿
      </Button>
      <Table dataSource={blocks} columns={columns} loading={loading} />
    </Card>
  );
}

export default BlockExplorer;