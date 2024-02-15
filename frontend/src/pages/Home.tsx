import React, {useEffect, useState} from 'react';
import {Button, message, Space, Table} from 'antd';
import CreatePaper from "../components/CreatePaper";
import {deletePaper, queryPaper, queryPapers} from "../apis/paper";
import {queryFileURL} from "../apis/file";
import ViewTranslated from "../components/ViewTranslated";

interface PaperItem {
  id: string;
  createAt: string;
  status: React.ReactDOM;
  ops: React.ReactDOM
}

const columns = [
  {
    title: 'ID',
    dataIndex: 'id',
    key: 'id',
  },
  {
    title: '创建时间',
    dataIndex: 'createAt',
    key: 'createAt',
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
  },
  {
    title: '操作',
    dataIndex: 'ops',
    key: 'ops',
  },
];


function Home() {

  const [
    messageApi,
    contextHolder
  ] = message.useMessage();

  const [creatorOpened, setCreatorOpened] = useState(false)
  const [viewOpened, setViewOpened] = useState(false)
  const [viewPaperID, setViewPaperID] = useState('')

  const [dataSource, setDataSource] = useState<PaperItem[]>([])

  const refreshPaperList = async () => {
    const papers = await queryPapers()
    setDataSource(papers.map((paper) => {
      const createAt = new Date(paper.createAt * 1000)
      let status = <div style={{color: '#C5FF00'}}>文字识别中</div>
      if (paper.status === 1) {
        status = <div style={{color: '#2357D8'}}>翻译中</div>
      } else if (paper.status === 2) {
        status = <div style={{color: '#16EF47'}}>完成</div>
      } else if (paper.status === 3) {
        status = <div style={{color: '#FF0000'}}>失败</div>
      }
      return {
        id: paper.paperID,
        createAt: `${createAt.getFullYear()}/${createAt.getMonth()}/${createAt.getDay()} ${createAt.getHours()}:${createAt.getMinutes()}`,
        status: status,
        ops: (<Space>
          <Button size={'small'} type={"link"} onClick={() => {
            queryPaper(paper.paperID).then((p) => {
              queryFileURL(p.fileHash).then(({url}) => {
                const link = document.createElement('a');
                link.style.display = 'none';
                link.href = url;
                link.setAttribute('download', "doc.pdf");
                document.body.appendChild(link);
                link.click();
                document.body.removeChild(link);
              })
            })
          }}>下载PDF</Button>
          <Button disabled={paper.status !== 2} size={"small"} type={"link"} onClick={() => {
            setViewPaperID(paper.paperID)
            setViewOpened(true)
          }}>查看结果</Button>
          <Button size={"small"} type={"link"} danger={true} onClick={() => {
            deletePaper(paper.paperID)
            refreshPaperList()
            messageApi.success("删除成功")
          }}>删除</Button>
        </Space>)
      };
    }) as unknown as PaperItem[])
  }

  useEffect(() => {
    refreshPaperList()
    setInterval(() => {
      refreshPaperList()
    }, 2000)
  }, [])

  return (
    <div style={{
      width: "60vw",
      margin: "40px auto"
    }}>
      {contextHolder}
      <div style={{
        fontSize: "20px",
        fontWeight: "bold"
      }}>
        论文翻译系统
      </div>
      <div style={{
        marginTop: "20px"
      }}>
        <Button type="primary" onClick={() => {
          setCreatorOpened(true)
        }}>创建</Button>
      </div>

      <div style={{
        marginTop: "20px"
      }}>
        <Table dataSource={dataSource} columns={columns}/>;
      </div>

      <CreatePaper opened={creatorOpened} onOk={() => {
        setCreatorOpened(false)
      }} onCancel={() => {
        setCreatorOpened(false)
      }}/>

      <ViewTranslated opened={viewOpened} paperID={viewPaperID} onOk={() => {
        setViewOpened(false)
      }} onCancel={() => {
        setViewOpened(false)
      }}/>
    </div>
  );
}

export default Home;
