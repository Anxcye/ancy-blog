package com.anxcye.constants;

public class SystemConstants
{
    /**
     *  文章是草稿
     */
    public static final int ARTICLE_STATUS_DRAFT = 1;
    /**
     *  文章是正常分布状态
     */
    public static final int ARTICLE_STATUS_NORMAL = 0;

    /**
     *  状态正常
     */
    public static final String STATUS_NORMAL = "0";

    /**
     *  分类状态正常
     */
    public static final String CATEGORY_STATUS_NORMAL = STATUS_NORMAL;

    /**
     *  友链状态正常
     */
    public static final String LINK_STATUS_NORMAL = STATUS_NORMAL;

    /**
     *  评论状态正常
     */
    public static final String COMMENT_STATUS_NORMAL = STATUS_NORMAL;

    /**
     *  评论是根评论
     */
    public static final long COMMENT_IS_ROOT = -1;

    public static final String COMMENT_TYPE_ARTICLE = "0";

    public static final String COMMENT_TYPE_LINK = "1";

    public static final long SUPER_ADMIN_ID = 1L;
}