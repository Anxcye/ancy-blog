package anxcye.es;

import org.apache.http.HttpHost;
import org.elasticsearch.action.admin.indices.delete.DeleteIndexRequest;
import org.elasticsearch.client.RequestOptions;
import org.elasticsearch.client.RestClient;
import org.elasticsearch.client.RestHighLevelClient;
import org.elasticsearch.client.indices.CreateIndexRequest;
import org.elasticsearch.client.indices.GetIndexRequest;
import org.elasticsearch.common.xcontent.XContentType;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.io.IOException;

public class IndexTest {

    private RestHighLevelClient client;

    @BeforeEach
    void setUp() {
        this.client = new RestHighLevelClient(RestClient.builder(
                HttpHost.create("http://localhost:9200")
        ));
    }

    @Test
    void testConnect() {
        System.out.println(client);
    }

    @Test
    void testCreateIndex() throws IOException {
        // 1.创建Request对象
        CreateIndexRequest request = new CreateIndexRequest("items");
        // 2.准备请求参数
        request.source(MAPPING_TEMPLATE, XContentType.JSON);
        // 3.发送请求
        client.indices().create(request, RequestOptions.DEFAULT);
    }

    @Test
    void testGetIndex() throws IOException {
        // 1.创建Request对象
        GetIndexRequest request = new GetIndexRequest("items");
        // 2.发送请求
        Boolean result = client.indices().exists(request, RequestOptions.DEFAULT);
        // 3.解析响应
        System.out.println(result);
    }

    @Test
    void testDeleteIndex() throws IOException {
        // 1.创建Request对象
        DeleteIndexRequest request = new DeleteIndexRequest("items");
        // 2.发送请求
        client.indices().delete(request, RequestOptions.DEFAULT);
    }

    @AfterEach
    void tearDown() throws IOException {
        this.client.close();
    }
    static final String MAPPING_TEMPLATE = """
            {
              "mappings": {
                "properties": {
                  "id": {
                    "type": "keyword"
                  },
                  "name":{
                    "type": "text",
                    "analyzer": "ik_max_word"
                  },
                  "price":{
                    "type": "integer"
                  },
                  "stock":{
                    "type": "integer"
                  },
                  "image":{
                    "type": "keyword",
                    "index": false
                  },
                  "category":{
                    "type": "keyword"
                  },
                  "brand":{
                    "type": "keyword"
                  },
                  "sold":{
                    "type": "integer"
                  },
                  "commentCount":{
                    "type": "integer"
                  },
                  "isAD":{
                    "type": "boolean"
                  },
                  "updateTime":{
                    "type": "date"
                  }
                }
              }
            }""";

}
