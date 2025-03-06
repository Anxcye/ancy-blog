package anxcye.es;

import cn.hutool.core.bean.BeanUtil;
import cn.hutool.json.JSONUtil;
import com.anxcye.AncyAdminApplication;
import com.anxcye.domain.entity.Article;
import com.anxcye.domain.po.ArticleDoc;
import com.anxcye.service.ArticleService;
import org.apache.http.HttpHost;
import org.elasticsearch.action.index.IndexRequest;
import org.elasticsearch.client.RequestOptions;
import org.elasticsearch.client.RestClient;
import org.elasticsearch.client.RestHighLevelClient;
import org.elasticsearch.common.xcontent.XContentType;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.context.annotation.ComponentScan;

import java.io.IOException;

@SpringBootTest(classes = AncyAdminApplication.class, properties = "spring.profiles.active=dev")
@ComponentScan(basePackages = "com.anxcye")
public class DocTest {
    private RestHighLevelClient client;

    @Autowired
    ArticleService articleService;

    @BeforeEach
    void setUp() {
        this.client = new RestHighLevelClient(RestClient.builder(
                HttpHost.create("http://localhost:9200")
        ));
    }


    @AfterEach
    void tearDown() throws IOException {
        this.client.close();
    }

    @Test
    void testAddDocument() throws IOException {
        // 1.根据id查询商品数据
        Article article = articleService.getById(1);
        // 2.转换为文档类型
        ArticleDoc itemDoc = BeanUtil.copyProperties(article, ArticleDoc.class);
        // 3.将ItemDTO转json
        String doc = JSONUtil.toJsonStr(itemDoc);

        // 1.准备Request对象
        IndexRequest request = new IndexRequest("items").id(itemDoc.getId().toString());
        // 2.准备Json文档
        request.source(doc, XContentType.JSON);
        // 3.发送请求
        client.index(request, RequestOptions.DEFAULT);
    }
}
