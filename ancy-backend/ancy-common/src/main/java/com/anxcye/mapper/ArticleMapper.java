package com.anxcye.mapper;

import com.anxcye.domin.entity.Article;
import com.baomidou.mybatisplus.core.mapper.BaseMapper;
import org.apache.ibatis.annotations.Mapper;

/**
* @author axy
* @description 针对表【ancy_article(文章表)】的数据库操作Mapper
* @createDate 2024-09-05 11:29:50
* @Entity .entity.Article
*/
@Mapper
public interface ArticleMapper extends BaseMapper<Article> {


}
